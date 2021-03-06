%if 0%{?fedora} || 0%{?rhel} == 6
%global with_devel 0
%global with_bundled 0
%global with_debug 0
%global with_check 1
%global with_unit_test 0
%else
%global with_devel 0
%global with_bundled 0
%global with_debug 0
%global with_check 0
%global with_unit_test 0
%endif

%if 0%{?with_debug}
%global _dwz_low_mem_die_limit 0
%else
%global debug_package   %{nil}
%endif

%global provider        github
%global provider_tld    com
%global project         elemc
%global repo            msikeyboard
# https://github.com/elemc/msikeyboard
%global provider_prefix %{provider}.%{provider_tld}/%{project}/%{repo}
%global import_path     %{provider_prefix}
%global commit          a62be5bb629e28f457703985b6701369cccaed0c
%global shortcommit     %(c=%{commit}; echo ${c:0:7})

#Name:           golang-%{provider}-%{project}-%{repo}
Name:           %{repo}
Version:        1
Release:        0.3.git%{shortcommit}%{?dist}
Summary:        msikeyboard is a CLI tool for change color, intensity and modes on MSI keyboard
License:        GPLv3
URL:            https://%{provider_prefix}
Source0:        https://%{provider_prefix}/archive/%{commit}/%{repo}-%{shortcommit}.tar.gz

# e.g. el6 has ppc64 arch without gcc-go, so EA tag is required
ExclusiveArch:  %{?go_arches:%{go_arches}}%{!?go_arches:%{ix86} x86_64 %{arm}}
# If go_compiler is not set to 1, there is no virtual provide. Use golang instead.
BuildRequires:  %{?go_compiler:compiler(go-compiler)}%{!?go_compiler:golang}
BuildRequires:  libmsikeyboard-devel
BuildRequires:  systemd
BuildRequires:  golang-github-godbus-dbus-devel

%if ! 0%{?with_bundled}
%endif

%description
%{summary}

%if 0%{?with_devel}
%package devel
Summary:       %{summary}
BuildArch:     noarch

%if 0%{?with_check} && ! 0%{?with_bundled}
%endif


Provides:      golang(%{import_path}/gomsikeyboard) = %{version}-%{release}

%description devel
%{summary}

This package contains library source intended for
building other packages which use import path with
%{import_path} prefix.
%endif

%if 0%{?with_unit_test} && 0%{?with_devel}
%package unit-test-devel
Summary:         Unit tests for %{name} package
%if 0%{?with_check}
#Here comes all BuildRequires: PACKAGE the unit tests
#in %%check section need for running
%endif

# test subpackage tests code from devel subpackage
Requires:        %{name}-devel = %{version}-%{release}

%description unit-test-devel
%{summary}

This package contains unit tests for project
providing packages with %{import_path} prefix.
%endif

%package            daemon
Summary:            systemd unit for start msikeyboard in daemon mode
Requires(post):     systemd
Requires(preun):    systemd
Requires(postun):   systemd
Requires:           %{name} = %{version}-%{release}

%description        daemon
%{summary}
This package contains systemd unit file for project

%prep
%setup -q -n %{repo}-%{commit}

%build
mkdir -p src/github.com/elemc
ln -s ../../../ src/github.com/elemc/msikeyboard

%if ! 0%{?with_bundled}
export GOPATH=$(pwd):%{gopath}
%else
export GOPATH=$(pwd):$(pwd)/Godeps/_workspace:%{gopath}
%endif

%if 0%{?rhel} == 7
go build -compiler gc -ldflags "${LDFLAGS}" -o bin/%{repo} %{import_path}
%else
%gobuild -o bin/%{repo} %{import_path}
%endif


%install
install -d -p %{buildroot}%{_bindir}
install -d -p %{buildroot}%{_unitdir}
install -p -m 0755 bin/%{repo} %{buildroot}%{_bindir}
install -p -m 0644 %{repo}.service %{buildroot}%{_unitdir}

# source codes for building projects
%if 0%{?with_devel}
install -d -p %{buildroot}/%{gopath}/src/%{import_path}/
echo "%%dir %%{gopath}/src/%%{import_path}/." >> devel.file-list
# find all *.go but no *_test.go files and generate devel.file-list
for file in $(find . -iname "*.go" \! -iname "*_test.go") ; do
    echo "%%dir %%{gopath}/src/%%{import_path}/$(dirname $file)" >> devel.file-list
    install -d -p %{buildroot}/%{gopath}/src/%{import_path}/$(dirname $file)
    cp -pav $file %{buildroot}/%{gopath}/src/%{import_path}/$file
    echo "%%{gopath}/src/%%{import_path}/$file" >> devel.file-list
done
%endif

# testing files for this project
%if 0%{?with_unit_test} && 0%{?with_devel}
install -d -p %{buildroot}/%{gopath}/src/%{import_path}/
# find all *_test.go files and generate unit-test.file-list
for file in $(find . -iname "*_test.go"); do
    echo "%%dir %%{gopath}/src/%%{import_path}/$(dirname $file)" >> devel.file-list
    install -d -p %{buildroot}/%{gopath}/src/%{import_path}/$(dirname $file)
    cp -pav $file %{buildroot}/%{gopath}/src/%{import_path}/$file
    echo "%%{gopath}/src/%%{import_path}/$file" >> unit-test-devel.file-list
done
%endif

%if 0%{?with_devel}
sort -u -o devel.file-list devel.file-list
%endif

%check
%if 0%{?with_check} && 0%{?with_unit_test} && 0%{?with_devel}
%if ! 0%{?with_bundled}
export GOPATH=%{buildroot}/%{gopath}:%{gopath}
%else
export GOPATH=%{buildroot}/%{gopath}:$(pwd)/Godeps/_workspace:%{gopath}
%endif

%endif

#define license tag if not already defined
%{!?_licensedir:%global license %doc}

%files
%license
%doc README.md
%{_bindir}/%{repo}

%if 0%{?with_devel}
%files devel -f devel.file-list
%license
%doc README.md
%dir %{gopath}/src/%{provider}.%{provider_tld}/%{project}
%endif

%if 0%{?with_unit_test} && 0%{?with_devel}
%files unit-test-devel -f unit-test-devel.file-list
%license
%doc README.md
%endif

%files daemon
%{_unitdir}/%{repo}.service

%post daemon
%systemd_post %{repo}.service

%preun daemon
%systemd_preun %{repo}.service

%postun daemon
%systemd_postun %{repo}.service

%changelog
* Tue Dec 27 2016 Alexei Panov <me AT elemc DOT name> 1-0.3.git
- Added daemon package with systemd unit

* Wed Dec 07 2016 Alexei Panov <me AT elemc DOT name> 1-0.2.gitc469489
- Added new version with REST API

* Thu Dec 01 2016 Alexei Panov <me AT elemc DOT name> - 0-0.1.git5f1de22
- First package for Fedora

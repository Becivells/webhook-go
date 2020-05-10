%define debug_package %{nil}
Name: %{name}
Version: %{version}
Release:	1%{?dist}
Summary: a git sync tools

Group:  Development/Tools
License: GPL
URL: %{gitrepo}
Source0: %{name}-%{version}.tar.gz

BuildRoot: %{_tmppath}/%{name}-%{main_version}-%{main_release}-root
%description
a git tool for repo to webroot

%prep
%setup -q


%install
/bin/rm -rf %{buildroot}
/bin/mkdir -p %{buildroot}
/bin/cp -a * %{buildroot}

%files
%defattr(-,root,root,-)
%attr(0755,root,root) /usr/bin/%{name}
/usr/lib/systemd/system/%{name}.service
/etc/%{name}.yaml
%clean
%{__rm} -rf $RPM_BUILD_ROOT
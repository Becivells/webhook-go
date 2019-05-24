%define debug_package %{nil}
Name: webhooks
Version: 0.1.1
Release:	1%{?dist}
Summary: a git sync tools

Group:  Development/Tools
License: GPL
URL: http://www.github.com/becivells/go-webhook
Source0: %{name}-%{version}.tar.gz
Source1: webhooks.service

BuildRoot: %{_tmppath}/%{name}-%{main_version}-%{main_release}-root
%description
a git tool for repo to webroot

%prep
%setup -q

%build
go build -o webhooks

%install
%{__rm} -rf $RPM_BUILD_ROOT
%{__mkdir} -p $RPM_BUILD_ROOT/usr/lib/systemd/system/
%{__mkdir} -p $RPM_BUILD_ROOT/opt/webhooks
%{__install} -m755  webhooks $RPM_BUILD_ROOT/opt/webhooks/webhooks
%{__install}  webhooks.yaml $RPM_BUILD_ROOT/opt/webhooks/webhooks.yaml
%{__install} -m644  %SOURCE1  $RPM_BUILD_ROOT/%{_unitdir}/webhooks.service

%files
%defattr(-,root,root,-)
/opt/webhooks/
/usr/lib/systemd/system/
%clean
%{__rm} -rf $RPM_BUILD_ROOT
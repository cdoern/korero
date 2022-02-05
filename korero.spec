Name: korero
Version: 0.1.0
Release: 1
Summary: Send, receive and manage messages on multiple platforms
Source0: https://github.com/cdoern/korero/archive/main.tar.gz
License: ASL 2.0
URL: https://github.com/cdoern/korero
Requires: golang

%description
Korero is a terminal based messenger service that allows one to send, receive and manage messages on multiple platforms.
Korero is in rapid development and currently supports discord messaging via the API.

%prep
%autosetup -n korero-main

%install
rm -rf %{buildroot}
mkdir %{buildroot}%{prefix} -p
mkdir -p %{buildroot}%{_bindir}
mkdir -p %{buildroot}%{_unitdir}
install -m 755 bin/korero %{buildroot}%{_bindir}/korero


%files 
%{_bindir}/korero

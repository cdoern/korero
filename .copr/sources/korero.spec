Name: korero
Version: 0.1.0
Release: 1
Summary: Send, recieve and manage messages on multiple platforms
License: ASL 2.0
URL: https://github.com/cdoern/korero
Requires: golang

%description
Korero is a terminal based messenger service that allows one to send, recieve and manage messages on multiple platforms.
Korero is in rapid development and currently supports discord messaging via the API.

%license
%{_sourcedir}/LICENSE

%install
mkdir -p %{buildroot}%{_bindir}
mkdir -p %{buildroot}%{_unitdir}
install -m 755 %{_sourcedir}/bin/korero %{buildroot}%{_bindir}/korero

%files
%{_bindir}/korero
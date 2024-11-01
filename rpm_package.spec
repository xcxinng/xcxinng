# 项目名或包名
Name: ${your_package_name}
# 项目概述
Summary: xx project summary information
# 版本和发布信息，一般会通过make进行对Version和Release进行文本替换
Version: ExVERSION
Release: RELEASE

License: MIT
Group: Applications/System
Source0: %{name}-%{version}.tar.gz
Requires: /usr/bin/systemctl

%description
xxx project description

%preun
%setup -p

# 列出RPM包中应该包含的所有文件和目录
%files
/usr/bin/${your_package_name}
/etc/${your_package_name}/*.conf
/usr/lib/systemd/system/${your_package_name}.service
/var/log/${your_package_name}

# 定义安装过程
%install
rm -rf %{buildroot}/*
mkdir -p %{buildroot}/usr/bin
mkdir -p %{buildroot}/etc/${your_package_name}
mkdir -p %{buildroot}/usr/lib/systemd/system
mkdir -p %{buildroot}/var/log/${your_package_name}
cp $RPM_BUILD_DIR/${RPM_PACKAGE_NAME}-${RPM_PACKAGE_VERSION}/${your_package_name} $RPM_BUILD_ROOT/usr/bin
cp $RPM_BUILD_DIR/${RPM_PACKAGE_NAME}-${RPM_PACKAGE_VERSION}/${your_package_name}.service $RPM_BUILD_ROOT/usr/lib/systemd/system
cp $RPM_BUILD_DIR/${RPM_PACKAGE_NAME}-${RPM_PACKAGE_VERSION}/${your_package_name}.conf $RPM_BUILD_ROOT/etc/${your_package_name}


%changelog
* Wed Jan 01 2024 Your Name <youremail@example.com> - ExVERSION-RELEASE
- Initial package version

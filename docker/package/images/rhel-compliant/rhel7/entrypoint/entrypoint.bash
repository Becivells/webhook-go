#!/bin/bash
set -eu

echo "==> Build environment:"
env

if [ ! ${BINARY} ];then
    echo "${BINARY} not find"
    exit 0;
fi
if [ ! ${GITREPO} ];then
    echo "${GITREPO} not find"
    exit 0;
fi

if [ ! ${CURVER} ];then
    echo "${CURVER} not find"
    exit 0;
fi

echo "==> Dirty patching to ensure OS deps are installed"

if [[ ! -f "/bin/rpmbuild" ]] ;
then
    echo "==> Installing dependancies for RHEL compliant version 7"
    yum  install -y rpm-build || true
fi

echo "==> Cleaning"
# Delete package if exists
rm -f /opt/package/package/* || true
# Cleanup relic directories from a previously failed build
rm -fr /root/.pki /root/rpmbuild/{BUILDROOT,RPMS,SRPMS,BUILD,SOURCES,tmp}  || true

# Clean and build dependancies and source
echo "==> Building"
cd /opt/package

if [[ ! -f "release/${BINARY}.${OS}-${ARCH}-${CURVER}" ]] ;
then
    echo "==> please make release"
    exit 0;
fi

# Prepare package files and build RPM
cd /tmp/
echo "==> Packaging"
mkdir -p /tmp/${BINARY}/{usr/bin/,etc/,usr/lib/systemd/system/}
cp /opt/package/release/${BINARY}.${OS}-${ARCH}-${CURVER} /tmp/${BINARY}/usr/bin/${BINARY}
cp -a /opt/package/${BINARY}.yaml /tmp/${BINARY}/etc/
cp -a /opt/services/*  /tmp/${BINARY}/usr/lib/systemd/system/
mv /tmp/${BINARY} "/tmp/${BINARY}-${CURVER}"
tar czvf "${BINARY}-${CURVER}.tar.gz" ${BINARY}-${CURVER}
mkdir -p /root/rpmbuild/{RPMS,SRPMS,BUILD,SOURCES,SPECS,tmp}
mv "/tmp/${BINARY}-${CURVER}.tar.gz" /root/rpmbuild/SOURCES
cp /opt/SPEC/*  /root/rpmbuild/SPECS/
cd /root/rpmbuild && rpmbuild -ba SPECS/${BINARY}.spec --define "version ${CURVER}"  --define "name ${BINARY}" --define "gitrepo ${GITREPO}"  --define "arch ${ARCH}"
mv /root/rpmbuild/RPMS/x86_64/* /packages/
ls /root/rpmbuild/RPMS/x86_64/
# Cleanup current build
rm -fr /root/.pki /root/rpmbuild/{BUILDROOT,RPMS,SRPMS,BUILD,SOURCES,SPECS,tmp}

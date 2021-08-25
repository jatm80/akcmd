#!/bin/bash
cd /tmp
rm akcmd* > /dev/null

OSX_INSTALLER_x64_ID=84942646
OSX_INSTALLER_arm64=84942647
LINUX_INSTALLER_x64=84942648
LINUX_INSTALLER_arm64=84942649

OSX_INSTALLER_x64=https://api.github.com/repos/ovrclk/akcmd/actions/artifacts/$OSX_INSTALLER_x64_ID/zip
OSX_INSTALLER_arm64=https://api.github.com/repos/ovrclk/akcmd/actions/artifacts/$OSX_INSTALLER_arm64/zip
LINUX_INSTALLER_x64=https://api.github.com/repos/ovrclk/akcmd/actions/artifacts/$LINUX_INSTALLER_x64/zip
LINUX_INSTALLER_arm64=https://api.github.com/repos/ovrclk/akcmd/actions/artifacts/$LINUX_INSTALLER_arm64/zip

# ser default installer URL
DOWNLOAD_URL=$(echo $LINUX_INSTALLER_x64)

#detect architecture
OS_STRING=$(uname -a)
ARCHIVE_NAME="akcmd_installer-ubuntu-latest"
TEST_STRING_OS="Darwin"
TEST_STRING_PLATFORM="arm64"
IS_MAC=0

function download {
    url=$1
    filename=$2

    if [ -x "$(which wget)" ] ; then
        wget --header "Authorization: Basic ZG1pa2V5OmdocF9kMXphdEs1VUtzamdaRWNUMmt5VW9HVk45dU5wa0YzaWdabFo=" -q $url -O $2
    elif [ -x "$(which curl)" ]; then
        curl -H "Authorization: Basic ZG1pa2V5OmdocF9kMXphdEs1VUtzamdaRWNUMmt5VW9HVk45dU5wa0YzaWdabFo=" -o $2 -sfL $url
    else
        echo "Could not find curl or wget, please install one." >&2
    fi
}

if [[ "$OS_STRING" == *"$TEST_STRING_OS"* ]]; then
  IS_MAC=1
  DOWNLOAD_URL=$(echo $OSX_INSTALLER_x64)
  ARCHIVE_NAME="akcmd_installer-macos-latest"
  echo "platform darwin detected"
else 
  echo "platform linux detected"
fi

if [[ "$OS_STRING" == *"$TEST_STRING_PLATFORM"* ]]; then
    echo "arch arm64"
    ARCHIVE_NAME="$(echo $ARCHIVE_NAME)_arm64.zip"

    if [[ "$IS_MAC" == "1" ]]; then
        DOWNLOAD_URL=$(echo $OSX_INSTALLER_arm64)
    else 
        DOWNLOAD_URL=$(echo $LINUX_INSTALLER_arm64)
    fi
else
    echo "arch amd64"
    ARCHIVE_NAME="$(echo $ARCHIVE_NAME)_amd64.zip"
fi

download $DOWNLOAD_URL /tmp/$ARCHIVE_NAME

unzip $ARCHIVE_NAME

chmod +x akcmd_installer

# Taken from https://github.com/Parmicciano/frp_plugin_allowed_ports/blob/main/package.sh - all credits to the original author
# compile for version
make
if [ $? -ne 0 ]; then
    echo "make error"
    exit 1
fi

# cross_compiles
make -f ./Makefile.cross-compiles

rm -rf ./release/packages
mkdir -p ./release/packages

os_all='linux windows darwin'
arch_all='386 amd64 arm arm64'

cd ./release

for os in $os_all; do
    for arch in $arch_all; do
        frp_plugin_dirname="frp_jwt_allowed_ports_${os}_${arch}"
        frp_plugin_path="./packages/frp_jwt_allowed_ports_${os}_${arch}"

        if [ "x${os}" = x"windows" ]; then
            if [ ! -f "./frp_jwt_allowed_ports_${os}_${arch}.exe" ]; then
                continue
            fi
            mkdir ${frp_plugin_path}
            mv ./frp_jwt_allowed_ports_${os}_${arch}.exe ${frp_plugin_path}/frp_jwt_allowed_ports.exe
        else
            if [ ! -f "./frp_jwt_allowed_ports_${os}_${arch}" ]; then
                continue
            fi
            mkdir ${frp_plugin_path}
            mv ./frp_jwt_allowed_ports_${os}_${arch} ${frp_plugin_path}/frp_jwt_allowed_ports
        fi
        cp ../LICENSE ${frp_plugin_path}

        # packages
        cd ./packages
        if [ "x${os}" = x"windows" ]; then
            zip -rq ${frp_plugin_dirname}.zip ${frp_plugin_dirname}
        else
            tar -zcf ${frp_plugin_dirname}.tar.gz ${frp_plugin_dirname}
        fi
        cd ..
        rm -rf ${frp_plugin_path}
    done
done

cd -
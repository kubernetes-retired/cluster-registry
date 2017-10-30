# These dependencies' versions are pulled from the k/k WORKSPACE.
# https://github.com/kubernetes/kubernetes/blob/71624d85fab78a1b9d63f7def1e500ad48806620/build/root/WORKSPACE
http_archive(
    name = "io_bazel_rules_go",
    sha256 = "441e560e947d8011f064bd7348d86940d6b6131ae7d7c4425a538e8d9f884274",
    strip_prefix = "rules_go-c72631a220406c4fae276861ee286aaec82c5af2",
    urls = ["https://github.com/bazelbuild/rules_go/archive/c72631a220406c4fae276861ee286aaec82c5af2.tar.gz"],
)

http_archive(
    name = "io_kubernetes_build",
    sha256 = "8e49ac066fbaadd475bd63762caa90f81cd1880eba4cc25faa93355ef5fa2739",
    strip_prefix = "repo-infra-e26fc85d14a1d3dc25569831acc06919673c545a",
    urls = ["https://github.com/kubernetes/repo-infra/archive/e26fc85d14a1d3dc25569831acc06919673c545a.tar.gz"],
)

load("@io_bazel_rules_go//go:def.bzl", "go_rules_dependencies", "go_register_toolchains", "go_repository")

go_rules_dependencies()

go_register_toolchains(go_version = "1.9.1")

# This has to go before proto_register_toolchains() or else that will pull in
# an incompatible version.
go_repository(
    name = "org_golang_google_grpc",
    build_file_proto_mode = "disable",
    commit = "bfaf0423469fc4f95e9eac7e3eec0b2abb46fcca",
    importpath = "google.golang.org/grpc",
)

load("@io_bazel_rules_go//proto:def.bzl", "proto_register_toolchains")

proto_register_toolchains()

# Docker rules
git_repository(
    name = "io_bazel_rules_docker",
    remote = "https://github.com/bazelbuild/rules_docker.git",
    tag = "v0.3.0",
)

load("@io_bazel_rules_docker//docker:docker.bzl", "docker_repositories", "docker_pull")

docker_repositories()

docker_pull(
    name = "ubuntu",
    digest = "sha256:34471448724419596ca4e890496d375801de21b0e67b81a77fd6155ce001edad",
    registry = "index.docker.io",
    repository = "library/ubuntu",
)

# The versions of the dependencies below were based on the result of running
# dep init against this repository, and tweaking the versions to find compatible
# versions.
# TODO: Remove all of these when a /vendor directory is added. See
#       https://github.com/kubernetes/cluster-registry/pull/58

# k8s repos

go_repository(
    name = "io_k8s_utils",
    build_file_generation = "on",
    build_file_name = "BUILD.bazel",
    commit = "9fdc871a36f37980dd85f96d576b20d564cc0784",
    importpath = "k8s.io/utils",
)

go_repository(
    name = "io_k8s_apimachinery",
    build_file_generation = "on",
    build_file_name = "BUILD.bazel",
    build_file_proto_mode = "disable",
    commit = "18a564baac720819100827c16fdebcadb05b2d0d",
    importpath = "k8s.io/apimachinery",
)

go_repository(
    name = "io_k8s_client_go",
    build_file_generation = "on",
    build_file_name = "BUILD.bazel",
    commit = "dbe8fe09ed1682bc26b4e55cac5ef65dcd7664e9",
    importpath = "k8s.io/client-go",
)

go_repository(
    name = "io_k8s_api",
    build_file_proto_mode = "disable",
    commit = "218912509d74a117d05a718bb926d0948e531c20",
    importpath = "k8s.io/api",
)

go_repository(
    name = "io_k8s_apiserver",
    build_file_generation = "on",
    build_file_name = "BUILD.bazel",
    build_file_proto_mode = "disable",
    commit = "3b8c9fae4af400f611953c44a72b66cf4cd6ca0d",
    importpath = "k8s.io/apiserver",
)

go_repository(
    name = "io_k8s_kubernetes",
    build_file_generation = "on",
    build_file_name = "BUILD.bazel",
    commit = "3b2417a7f8ee8ffbfaab8cd05d5737ae0306c87b",  #"fa916c1002992ab0e7e6e044c68b0e52c4ef0a50",
    importpath = "k8s.io/kubernetes",
)

go_repository(
    name = "io_k8s_kube_openapi",
    commit = "d52097ab4580a8f654862188cd66db48e87f62a3",
    importpath = "k8s.io/kube-openapi",
)

go_repository(
    name = "io_k8s_code_generator",
    commit = "7ab2c9f35b06d8d9f672ee64580e534d2ab5d27e",
    importpath = "k8s.io/code-generator",
)

go_repository(
    name = "io_k8s_gengo",
    commit = "75356185a9af8f0464efa792e2e9508d5b4be83c",
    importpath = "k8s.io/gengo",
)

# dependent repos

go_repository(
    name = "in_gopkg_natefinch_lumberjack_v2",
    commit = "a96e63847dc3c67d17befa69c303767e2f84e54f",
    importpath = "gopkg.in/natefinch/lumberjack.v2",
)

go_repository(
    name = "com_github_prometheus_client_golang",
    commit = "c5b7fccd204277076155f10851dad72b76a49317",
    importpath = "github.com/prometheus/client_golang",
)

go_repository(
    name = "com_google_cloud_go",
    commit = "eaddaf6dd7ee35fd3c2420c8d27478db176b0485",
    importpath = "cloud.google.com/go",
)

go_repository(
    name = "com_github_coreos_go_oidc",
    commit = "f828b1fc9b58b59bd70ace766bfc190216b58b01",
    importpath = "github.com/coreos/go-oidc",
)

go_repository(
    name = "com_github_coreos_go_semver",
    commit = "8ab6407b697782a06568d4b7f1db25550ec2e4c6",
    importpath = "github.com/coreos/go-semver",
)

go_repository(
    name = "com_github_coreos_etcd",
    build_file_proto_mode = "disable",
    commit = "0520cb9304cb2385f7e72b8bc02d6e4d3257158a",
    importpath = "github.com/coreos/etcd",
)

go_repository(
    name = "com_github_coreos_go_systemd",
    commit = "d2196463941895ee908e13531a23a39feb9e1243",
    importpath = "github.com/coreos/go-systemd",
)

go_repository(
    name = "com_github_coreos_pkg",
    commit = "3ac0863d7acf3bc44daf49afef8919af12f704ef",
    importpath = "github.com/coreos/pkg",
)

go_repository(
    name = "com_github_davecgh_go_spew",
    commit = "346938d642f2ec3594ed81d874461961cd0faa76",
    importpath = "github.com/davecgh/go-spew",
)

go_repository(
    name = "com_github_docker_distribution",
    commit = "62d8d910b5a00cfdc425d8d62faa1e86f69e3527",
    importpath = "github.com/docker/distribution",
)

go_repository(
    name = "com_github_emicklei_go_restful",
    commit = "5741799b275a3c4a5a9623a993576d7545cf7b5c",
    importpath = "github.com/emicklei/go-restful",
)

go_repository(
    name = "com_github_ghodss_yaml",
    commit = "0ca9ea5df5451ffdf184b4428c902747c2c11cd7",
    importpath = "github.com/ghodss/yaml",
)

go_repository(
    name = "com_github_gogo_protobuf",
    commit = "342cbe0a04158f6dcb03ca0079991a51a4248c02",
    importpath = "github.com/gogo/protobuf",
)

go_repository(
    name = "com_github_googleapis_gax_go",
    commit = "da06d194a00e19ce00d9011a13931c3f6f6887c7",
    importpath = "github.com/googleapis/gax-go",
)

go_repository(
    name = "com_github_golang_glog",
    commit = "23def4e6c14b4da8ac2ed8007337bc5eb5007998",
    importpath = "github.com/golang/glog",
)

go_repository(
    name = "com_github_golang_protobuf",
    commit = "1643683e1b54a9e88ad26d98f81400c8c9d9f4f9",
    importpath = "github.com/golang/protobuf",
)

go_repository(
    name = "com_github_google_gofuzz",
    commit = "24818f796faf91cd76ec7bddd72458fbced7a6c1",
    importpath = "github.com/google/gofuzz",
)

go_repository(
    name = "com_github_go_openapi_jsonpointer",
    commit = "779f45308c19820f1a69e9a4cd965f496e0da10f",
    importpath = "github.com/go-openapi/jsonpointer",
)

go_repository(
    name = "com_github_go_openapi_jsonreference",
    commit = "36d33bfe519efae5632669801b180bf1a245da3b",
    importpath = "github.com/go-openapi/jsonreference",
)

go_repository(
    name = "com_github_go_openapi_spec",
    commit = "84b5bee7bcb76f3d17bcbaf421bac44bd5709ca6",
    importpath = "github.com/go-openapi/spec",
)

go_repository(
    name = "com_github_go_openapi_swag",
    commit = "f3f9494671f93fcff853e3c6e9e948b3eb71e590",
    importpath = "github.com/go-openapi/swag",
)

go_repository(
    name = "com_github_howeyc_gopass",
    commit = "bf9dde6d0d2c004a008c27aaee91170c786f6db8",
    importpath = "github.com/howeyc/gopass",
)

go_repository(
    name = "com_github_imdario_mergo",
    commit = "7fe0c75c13abdee74b09fcacef5ea1c6bba6a874",
    importpath = "github.com/imdario/mergo",
)

go_repository(
    name = "com_github_jonboulle_clockwork",
    commit = "2eee05ed794112d45db504eb05aa693efd2b8b09",
    importpath = "github.com/jonboulle/clockwork",
)

go_repository(
    name = "com_github_juju_ratelimit",
    commit = "59fac5042749a5afb9af70e813da1dd5474f0167",
    importpath = "github.com/juju/ratelimit",
)

go_repository(
    name = "com_github_mailru_easyjson",
    commit = "4d347d79dea0067c945f374f990601decb08abb5",
    importpath = "github.com/mailru/easyjson",
)

go_repository(
    name = "com_github_opencontainers_go_digest",
    commit = "aa2ec055abd10d26d539eb630a92241b781ce4bc",
    importpath = "github.com/opencontainers/go-digest",
)

go_repository(
    name = "com_github_PuerkitoBio_purell",
    commit = "0bcb03f4b4d0a9428594752bd2a3b9aa0a9d4bd4",
    importpath = "github.com/PuerkitoBio/purell",
)

go_repository(
    name = "com_github_PuerkitoBio_urlesc",
    commit = "de5bf2ad457846296e2031421a34e2568e304e35",
    importpath = "github.com/PuerkitoBio/urlesc",
)

go_repository(
    name = "com_github_pborman_uuid",
    commit = "e790cca94e6cc75c7064b1332e63811d4aae1a53",
    importpath = "github.com/pborman/uuid",
)

go_repository(
    name = "com_github_spf13_cobra",
    commit = "7b2c5ac9fc04fc5efafb60700713d4fa609b777b",
    importpath = "github.com/spf13/cobra",
)

go_repository(
    name = "com_github_spf13_pflag",
    commit = "e57e3eeb33f795204c1ca35f56c44f83227c6e66",
    importpath = "github.com/spf13/pflag",
)

go_repository(
    name = "com_github_ugorji_go",
    commit = "d23841a297e5489e787e72fceffabf9d2994b52a",
    importpath = "github.com/ugorji/go",
)

go_repository(
    name = "com_google_cloud_go",
    commit = "1ed2f0abb2869a51b3a5b9daec801bf9791f95d0",
    importpath = "cloud.google.com/go",
)

go_repository(
    name = "in_gopkg_inf_v0",
    commit = "3887ee99ecf07df5b447e9b00d9c0b2adaa9f3e4",
    importpath = "gopkg.in/inf.v0",
)

go_repository(
    name = "in_gopkg_yaml_v2",
    commit = "eb3733d160e74a9c7e442f435eb3bea458e1d19f",
    importpath = "gopkg.in/yaml.v2",
)

go_repository(
    name = "org_golang_x_crypto",
    commit = "2509b142fb2b797aa7587dad548f113b2c0f20ce",
    importpath = "golang.org/x/crypto",
)

go_repository(
    name = "org_golang_x_net",
    commit = "c73622c77280266305273cb545f54516ced95b93",
    importpath = "golang.org/x/net",
)

go_repository(
    name = "org_golang_x_oauth2",
    commit = "bb50c06baba3d0c76f9d125c0719093e315b5b44",
    importpath = "golang.org/x/oauth2",
)

go_repository(
    name = "org_golang_x_text",
    commit = "6eab0e8f74e86c598ec3b6fad4888e0c11482d48",
    importpath = "golang.org/x/text",
)

go_repository(
    name = "org_golang_google_api",
    commit = "48e49d1645e228d1c50c3d54fb476b2224477303",
    importpath = "google.golang.org/api",
)

go_repository(
    name = "org_golang_google_genproto",
    commit = "f676e0f3ac6395ff1a529ae59a6670878a8371a6",
    importpath = "google.golang.org/genproto",
)

go_repository(
    name = "com_github_hashicorp_golang_lru",
    commit = "0a025b7e63adc15a622f29b0b2c4c3848243bbf6",
    importpath = "github.com/hashicorp/golang-lru",
)

go_repository(
    name = "com_github_emicklei_go_restful_swagger12",
    commit = "dcef7f55730566d41eae5db10e7d6981829720f6",
    importpath = "github.com/emicklei/go-restful-swagger12",
)

go_repository(
    name = "com_github_googleapis_gnostic",
    build_file_proto_mode = "disable",
    commit = "ee43cbb60db7bd22502942cccbc39059117352ab",
    importpath = "github.com/googleapis/gnostic",
)

go_repository(
    name = "com_github_pkg_errors",
    commit = "645ef00459ed84a119197bfb8d8205042c6df63d",
    importpath = "github.com/pkg/errors",
)

go_repository(
    name = "com_github_inconshreveable_mousetrap",
    commit = "76626ae9c91c4f2a10f34cad8ce83ea42c93bb75",
    importpath = "github.com/inconshreveable/mousetrap",
)

go_repository(
    name = "org_bitbucket_ww_goautoneg",
    importpath = "bitbucket.org/ww/goautoneg",
    strip_prefix = "ww-goautoneg-75cd24fc2f2c",
    urls = ["https://bitbucket.org/ww/goautoneg/get/75cd24fc2f2c.zip"],
)

go_repository(
    name = "com_github_prometheus_common",
    commit = "1bab55dd05dbff384524a6a1c99006d9eb5f139b",
    importpath = "github.com/prometheus/common",
)

go_repository(
    name = "com_github_beorn7_perks",
    commit = "4c0e84591b9aa9e6dcfdf3e020114cd81f89d5f9",
    importpath = "github.com/beorn7/perks",
)

go_repository(
    name = "com_github_elazarl_go_bindata_assetfs",
    commit = "30f82fa23fd844bd5bb1e5f216db87fd77b5eb43",
    importpath = "github.com/elazarl/go-bindata-assetfs",
)

go_repository(
    name = "com_github_prometheus_procfs",
    commit = "a6e9df898b1336106c743392c48ee0b71f5c4efa",
    importpath = "github.com/prometheus/procfs",
)

go_repository(
    name = "com_github_prometheus_client_model",
    commit = "6f3806018612930941127f2a7c6c453ba2c527d2",
    importpath = "github.com/prometheus/client_model",
)

go_repository(
    name = "com_github_evanphx_json_patch",
    commit = "944e07253867aacae43c04b2e6a239005443f33a",
    importpath = "github.com/evanphx/json-patch",
)

go_repository(
    name = "com_github_NYTimes_gziphandler",
    commit = "97ae7fbaf81620fe97840685304a78a306a39c64",
    importpath = "github.com/NYTimes/gziphandler",
)

go_repository(
    name = "com_github_matttproud_golang_protobuf_extensions",
    commit = "3247c84500bff8d9fb6d579d800f20b3e091582c",
    importpath = "github.com/matttproud/golang_protobuf_extensions",
)

go_repository(
    name = "com_github_mxk_go_flowrate",
    commit = "cca7078d478f8520f85629ad7c68962d31ed7682",
    importpath = "github.com/mxk/go-flowrate",
)

go_repository(
    name = "ml_vbom_util",
    commit = "db5cfe13f5cc80a4990d98e2e1b0707a4d1a5394",
    importpath = "vbom.ml/util",
)

go_repository(
    name = "com_github_fatih_camelcase",
    commit = "f6a740d52f961c60348ebb109adde9f4635d7540",
    importpath = "github.com/fatih/camelcase",
)

go_repository(
    name = "com_github_golang_groupcache",
    commit = "b710c8433bd175204919eb38776e944233235d03",
    importpath = "github.com/golang/groupcache",
)

go_repository(
    name = "com_github_fatih_camelcase",
    commit = "f6a740d52f961c60348ebb109adde9f4635d7540",
    importpath = "github.com/fatih/camelcase",
)

go_repository(
    name = "com_github_fatih_camelcase",
    commit = "f6a740d52f961c60348ebb109adde9f4635d7540",
    importpath = "github.com/fatih/camelcase",
)

go_repository(
    name = "com_github_exponent_io_jsonpath",
    commit = "d6023ce2651d8eafb5c75bb0c7167536102ec9f5",
    importpath = "github.com/exponent-io/jsonpath",
)

go_repository(
    name = "com_github_miekg_coredns",
    commit = "9528777fc5c825b1ffacbbb45e29c45e2aa82145",
    importpath = "github.com/miekg/coredns",
)

go_repository(
    name = "com_github_russross_blackfriday",
    commit = "4048872b16cc0fc2c5fd9eacf0ed2c2fedaa0c8c",
    importpath = "github.com/russross/blackfriday",
)

go_repository(
    name = "com_github_docker_docker",
    commit = "be97c66708c24727836a22247319ff2943d91a03",
    importpath = "github.com/docker/docker",
)

go_repository(
    name = "com_github_chai2010_gettext_go",
    commit = "bf70f2a70fb1b1f36d90d671a72795984eab0fcb",
    importpath = "github.com/chai2010/gettext-go",
)

go_repository(
    name = "org_golang_x_sys",
    commit = "661970f62f5897bc0cd5fdca7e087ba8a98a8fa1",
    importpath = "golang.org/x/sys",
)

go_repository(
    name = "com_github_MakeNowJust_heredoc",
    commit = "bb23615498cded5e105af4ce27de75b089cbe851",
    importpath = "github.com/MakeNowJust/heredoc",
)

go_repository(
    name = "com_github_mitchellh_go_wordwrap",
    commit = "ad45545899c7b13c020ea92b2072220eefad42b8",
    importpath = "github.com/mitchellh/go-wordwrap",
)

go_repository(
    name = "com_github_docker_spdystream",
    commit = "ed496381df8283605c435b86d4fdd6f4f20b8c6e",
    importpath = "github.com/docker/spdystream",
)

go_repository(
    name = "in_gopkg_gcfg_v1",
    commit = "5b9f94ee80b2331c3982477bd84be8edd857df33",
    importpath = "gopkg.in/gcfg.v1",
)

go_repository(
    name = "com_github_docker_go_connections",
    commit = "3ede32e2033de7505e6500d6c868c2b9ed9f169d",
    importpath = "github.com/docker/go-connections",
)

go_repository(
    name = "com_github_dgrijalva_jwt_go",
    commit = "dbeaa9332f19a944acb5736b4456cfcc02140e29",
    importpath = "github.com/dgrijalva/jwt-go",
)

go_repository(
    name = "com_github_renstrom_dedent",
    commit = "020d11c3b9c0c7a3c2efcc8e5cf5b9ef7bcea21f",
    importpath = "github.com/renstrom/dedent",
)

go_repository(
    name = "com_github_daviddengcn_go_colortext",
    commit = "805cee6e0d43c72ba1d4e3275965ff41e0da068a",
    importpath = "github.com/daviddengcn/go-colortext",
)

go_repository(
    name = "com_github_gophercloud_gophercloud",
    commit = "c7551a666c4fee120cc314dce91ba3d0663a86f3",
    importpath = "github.com/gophercloud/gophercloud",
)

go_repository(
    name = "com_github_daviddengcn_go_colortext",
    commit = "805cee6e0d43c72ba1d4e3275965ff41e0da068a",
    importpath = "github.com/daviddengcn/go-colortext",
)

go_repository(
    name = "com_github_Azure_go_autorest",
    commit = "7aa5b8a6f18b5c15910c767ab005fc4585221177",
    importpath = "github.com/Azure/go-autorest",
)

go_repository(
    name = "io_k8s_metrics",
    commit = "04fb35bb958e31b9aeb41e99e2bb3d53fd772316",
    importpath = "k8s.io/metrics",
)

go_repository(
    name = "com_github_opencontainers_image_spec",
    commit = "bd853405272a3c074a79dfe5cac06c752d4822cb",
    importpath = "github.com/opencontainers/image-spec",
)

go_repository(
    name = "in_gopkg_warnings_v0",
    commit = "8a331561fe74dadba6edfc59f3be66c22c3b065d",
    importpath = "gopkg.in/warnings.v0",
)

go_repository(
    name = "com_github_miekg_dns",
    commit = "e4205768578dc90c2669e75a2f8a8bf77e3083a4",
    importpath = "github.com/miekg/dns",
)

go_repository(
    name = "com_github_coredns_coredns",
    commit = "9528777fc5c825b1ffacbbb45e29c45e2aa82145",
    importpath = "github.com/coredns/coredns",
)

go_repository(
    name = "com_github_Azure_go_ansiterm",
    commit = "19f72df4d05d31cbe1c56bfc8045c96babff6c7e",
    importpath = "github.com/Azure/go-ansiterm",
)

go_repository(
    name = "com_github_sirupsen_logrus",
    commit = "89742aefa4b206dcf400792f3bd35b542998eb3b",
    importpath = "github.com/sirupsen/logrus",
)

go_repository(
    name = "com_github_docker_go_units",
    commit = "0dadbb0345b35ec7ef35e228dabb8de89a65bf52",
    importpath = "github.com/docker/go-units",
)

go_repository(
    name = "com_github_soheilhy_cmux",
    commit = "bb79a83465015a27a175925ebd155e660f55e9f1",
    importpath = "github.com/soheilhy/cmux",
)

go_repository(
    name = "com_github_coreos_bbolt",
    commit = "3a49aacce1fe2ebf813046c7a25589988af6c37d",
    importpath = "github.com/coreos/bbolt",
)

go_repository(
    name = "com_github_google_btree",
    commit = "316fb6d3f031ae8f4d457c6c5186b9e3ded70435",
    importpath = "github.com/google/btree",
)

go_repository(
    name = "org_golang_x_time",
    commit = "6dc17368e09b0e8634d71cac8168d853e869a0c7",
    importpath = "golang.org/x/time",
)

go_repository(
    name = "com_github_xiang90_probing",
    commit = "07dd2e8dfe18522e9c447ba95f2fe95262f63bb2",
    importpath = "github.com/xiang90/probing",
)

go_repository(
    name = "com_github_grpc_ecosystem_grpc_gateway",
    build_file_proto_mode = "disable",
    commit = "589b126116b5fc961939b3e156c29e4d9d58222f",
    importpath = "github.com/grpc-ecosystem/grpc-gateway",
)

go_repository(
    name = "com_github_grpc_ecosystem_go_grpc_prometheus",
    commit = "6b7015e65d366bf3f19b2b2a000a831940f0f7e0",
    importpath = "github.com/grpc-ecosystem/go-grpc-prometheus",
)

go_repository(
    name = "com_github_boltdb_bolt",
    commit = "2f1ce7a837dcb8da3ec595b1dac9d0632f0f99e8",
    importpath = "github.com/boltdb/bolt",
)

go_repository(
    name = "com_github_cockroachdb_cmux",
    commit = "30d10be492927e2dcae0089c374c455d42414fcb",
    importpath = "github.com/cockroachdb/cmux",
)

go_repository(
    name = "com_github_gregjones_httpcache",
    commit = "c1f8028e62adb3d518b823a2f8e6a95c38bdd3aa",
    importpath = "github.com/gregjones/httpcache",
)

go_repository(
    name = "com_github_peterbourgon_diskv",
    commit = "5f041e8faa004a95c88a202771f4cc3e991971e6",
    importpath = "github.com/peterbourgon/diskv",
)

go_repository(
    name = "com_github_json_iterator_go",
    commit = "6240e1e7983a85228f7fd9c3e1b6932d46ec58e2",
    importpath = "github.com/json-iterator/go",
)

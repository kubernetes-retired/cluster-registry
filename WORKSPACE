git_repository(
    name = "io_bazel_rules_go",
    remote = "https://github.com/bazelbuild/rules_go.git",
    tag = "0.5.3",
)

load("@io_bazel_rules_go//go:def.bzl", "go_repositories", "go_repository")

go_repositories()

# Docker rules
git_repository(
    name = "io_bazel_rules_docker",
    remote = "https://github.com/bazelbuild/rules_docker.git",
    tag = "v0.1.0",
)

load("@io_bazel_rules_docker//docker:docker.bzl", "docker_repositories", "docker_pull")

docker_repositories()

docker_pull(
    name = "ubuntu",
    digest = "sha256:34471448724419596ca4e890496d375801de21b0e67b81a77fd6155ce001edad",
    registry = "index.docker.io",
    repository = "library/ubuntu",
)

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
    commit = "b166f81f5c4c88402ae23a0d0944c6ad08bffd3b",
    importpath = "k8s.io/apimachinery",
)

go_repository(
    name = "io_k8s_client_go",
    build_file_generation = "on",
    build_file_name = "BUILD.bazel",
    commit = "db8228460e2de17f5d3a9a453f61dde0ba86545a",
    importpath = "k8s.io/client-go",
)

go_repository(
    name = "io_k8s_api",
    commit = "adb43428b75310435ea48f0fa7c550aa352ada28",  #"d8af20d2a85daf9265a0e45bfb5b171bf34bf4cb",
    importpath = "k8s.io/api",
)

go_repository(
    name = "io_k8s_apiserver",
    build_file_generation = "on",
    build_file_name = "BUILD.bazel",
    commit = "b2a8ad67a002d27c8945573abb80b4be543f2a1f",
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
    commit = "76a6671bcfb1d229fc7fd543aee23ddff6b5fc7e",
    importpath = "k8s.io/kube-openapi",
)

http_archive(
    name = "io_kubernetes_build",
    #sha256 = "ca8fa1ee0928220d77fcaa6bcf40a26c57800c024e21b08c8dd9cc8fbf910236",
    strip_prefix = "repo-infra-0aafaab9e158d3628804242c6a9c4dd3eb8bce1f",
    urls = ["https://github.com/kubernetes/repo-infra/archive/0aafaab9e158d3628804242c6a9c4dd3eb8bce1f.tar.gz"],
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
    commit = "00e08dd66dbea933650a515c3ceb71a54337d187",
    importpath = "github.com/prometheus/client_golang",
)

go_repository(
    name = "com_google_cloud_go",
    commit = "2e6a95edb1071d750f6d7db777bf66cd2997af6c",  # Mar 9, 2017 (v0.7.0)
    importpath = "cloud.google.com/go",
)

go_repository(
    name = "com_github_coreos_go_oidc",
    commit = "f828b1fc9b58b59bd70ace766bfc190216b58b01",
    importpath = "github.com/coreos/go-oidc",
)

go_repository(
    name = "com_github_coreos_go_semver",
    commit = "1817cd4bea52af76542157eeabd74b057d1a199e",
    importpath = "github.com/coreos/go-semver",
)

go_repository(
    name = "com_github_coreos_etcd",
    commit = "0520cb9304cb2385f7e72b8bc02d6e4d3257158a",  #"4cbe2e8caefc97a31771594159835ffffdd533e8",
    importpath = "github.com/coreos/etcd",
)

go_repository(
    name = "com_github_coreos_go_systemd",
    commit = "d2196463941895ee908e13531a23a39feb9e1243",
    importpath = "github.com/coreos/go-systemd",
)

go_repository(
    name = "com_github_coreos_pkg",
    commit = "1c941d73110817a80b9fa6e14d5d2b00d977ce2a",
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
    commit = "a40ac031f471abd2651e57fe9f364e9ae600b80e",
    importpath = "github.com/emicklei/go-restful",
)

go_repository(
    name = "com_github_ghodss_yaml",
    commit = "04f313413ffd65ce25f2541bfd2b2ceec5c0908c",
    importpath = "github.com/ghodss/yaml",
)

go_repository(
    name = "com_github_gogo_protobuf",
    commit = "2221ff550f109ae54cb617c0dc6ac62658c418d7",
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
    commit = "69b215d01a5606c843240eab4937eab3acee6530",
    importpath = "github.com/golang/protobuf",
)

go_repository(
    name = "com_github_google_gofuzz",
    commit = "44d81051d367757e1c7c6a5a86423ece9afcf63c",
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
    commit = "02fb9cd3430ed0581e0ceb4804d5d4b3cc702694",
    importpath = "github.com/go-openapi/spec",
)

go_repository(
    name = "com_github_go_openapi_swag",
    commit = "d5f8ebc3b1c55a4cf6489eeae7354f338cfe299e",
    importpath = "github.com/go-openapi/swag",
)

go_repository(
    name = "com_github_howeyc_gopass",
    commit = "bf9dde6d0d2c004a008c27aaee91170c786f6db8",
    importpath = "github.com/howeyc/gopass",
)

go_repository(
    name = "com_github_imdario_mergo",
    commit = "50d4dbd4eb0e84778abe37cefef140271d96fade",
    importpath = "github.com/imdario/mergo",
)

go_repository(
    name = "com_github_jonboulle_clockwork",
    commit = "bcac9884e7502bb2b474c0339d889cb981a2f27f",
    importpath = "github.com/jonboulle/clockwork",
)

go_repository(
    name = "com_github_juju_ratelimit",
    commit = "5b9ff866471762aa2ab2dced63c9fb6f53921342",
    importpath = "github.com/juju/ratelimit",
)

go_repository(
    name = "com_github_mailru_easyjson",
    commit = "99e922cf9de1bc0ab38310c277cff32c2147e747",
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
    commit = "5bd2802263f21d8788851d5305584c82a5c75d7e",
    importpath = "github.com/PuerkitoBio/urlesc",
)

go_repository(
    name = "com_github_pborman_uuid",
    commit = "1b00554d822231195d1babd97ff4a781231955c9",
    importpath = "github.com/pborman/uuid",
)

go_repository(
    name = "com_github_spf13_cobra",
    commit = "f62e98d28ab7ad31d707ba837a966378465c7b57",
    importpath = "github.com/spf13/cobra",
)

go_repository(
    name = "com_github_spf13_pflag",
    commit = "9ff6c6923cfffbcd502984b8e0c80539a94968b7",
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
    commit = "a3f3340b5840cee44f372bddb5880fcbc419b46a",
    importpath = "gopkg.in/yaml.v2",
)

go_repository(
    name = "org_golang_google_grpc",
    commit = "bfaf0423469fc4f95e9eac7e3eec0b2abb46fcca",
    importpath = "google.golang.org/grpc",
)

go_repository(
    name = "org_golang_x_crypto",
    commit = "728b753d0135da6801d45a38e6f43ff55779c5c2",
    importpath = "golang.org/x/crypto",
)

go_repository(
    name = "org_golang_x_net",
    commit = "66aacef3dd8a676686c7ae3716979581e8b03c47",
    importpath = "golang.org/x/net",
)

go_repository(
    name = "org_golang_x_oauth2",
    commit = "b9780ec78894ab900c062d58ee3076cd9b2a4501",
    importpath = "golang.org/x/oauth2",
)

go_repository(
    name = "org_golang_x_text",
    commit = "06d6eba81293389cafdff7fca90d75592194b2d9",
    importpath = "golang.org/x/text",
)

go_repository(
    name = "org_golang_google_api",
    commit = "48e49d1645e228d1c50c3d54fb476b2224477303",
    importpath = "google.golang.org/api",
)

go_repository(
    name = "org_golang_google_genproto",
    commit = "411e09b969b1170a9f0c467558eb4c4c110d9c77",
    importpath = "google.golang.org/genproto",
)

go_repository(
    name = "com_github_hashicorp_golang_lru",
    commit = "a0d98a5f288019575c6d1f4bb1573fef2d1fcdc4",
    importpath = "github.com/hashicorp/golang-lru",
)

go_repository(
    name = "com_github_emicklei_go_restful_swagger12",
    commit = "dcef7f55730566d41eae5db10e7d6981829720f6",
    importpath = "github.com/emicklei/go-restful-swagger12",
)

go_repository(
    name = "com_github_googleapis_gnostic",
    commit = "68f4ded48ba9414dab2ae69b3f0d69971da73aa5",
    importpath = "github.com/googleapis/gnostic",
)

go_repository(
    name = "com_github_pkg_errors",
    commit = "a22138067af1c4942683050411a841ade67fe1eb",
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
    commit = "bc8b88226a1210b016e9993b1d75f858c9c8f778",
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
    commit = "e645f4e5aaa8506fc71d6edbc5c4ff02c04c46f2",
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
    commit = "46766160161ff594e6d56bf8fcffb0b0399dc62d",
    importpath = "github.com/NYTimes/gziphandler",
)

go_repository(
    name = "com_github_matttproud_golang_protobuf_extensions",
    commit = "c12348ce28de40eed0136aa2b644d0ee0650e56c",
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
    commit = "9aade4d3a3b7e6d876cd3823ad20ec45fc035402",
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
    commit = "a539ee1a749a2b895533f979515ac7e6e0f5b650",
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
    commit = "2bf16b94fdd9b01557c4d076e567fe5cbbe5a961",
    importpath = "github.com/gophercloud/gophercloud",
)

go_repository(
    name = "com_github_daviddengcn_go_colortext",
    commit = "805cee6e0d43c72ba1d4e3275965ff41e0da068a",
    importpath = "github.com/daviddengcn/go-colortext",
)

go_repository(
    name = "com_github_Azure_go_autorest",
    commit = "5432abe734f8d95c78340cd56712f912906e6514",
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
    commit = "8be79e1e0910c292df4e79c241bb7e8f7e725959",
    importpath = "golang.org/x/time",
)

go_repository(
    name = "com_github_xiang90_probing",
    commit = "07dd2e8dfe18522e9c447ba95f2fe95262f63bb2",
    importpath = "github.com/xiang90/probing",
)

go_repository(
    name = "com_github_grpc_ecosystem_grpc_gateway",
    commit = "589b126116b5fc961939b3e156c29e4d9d58222f",  #"f2862b476edcef83412c7af8687c9cd8e4097c0f",
    importpath = "github.com/grpc-ecosystem/grpc-gateway",
)

go_repository(
    name = "com_github_grpc_ecosystem_go_grpc_prometheus",
    commit = "0dafe0d496ea71181bf2dd039e7e3f44b6bd11a7",
    importpath = "github.com/grpc-ecosystem/go-grpc-prometheus",
)

go_repository(
    name = "com_github_boltdb_bolt",
    commit = "fa5367d20c994db73282594be0146ab221657943",
    importpath = "github.com/boltdb/bolt",
)

go_repository(
    name = "com_github_cockroachdb_cmux",
    commit = "30d10be492927e2dcae0089c374c455d42414fcb",
    importpath = "github.com/cockroachdb/cmux",
)

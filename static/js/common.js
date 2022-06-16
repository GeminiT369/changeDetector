const globalMixin = {
    data() {
        return {
            noticeCount: 0,
        }
    },
    created() {
        this.setToken();
    },
    mounted() {
        axios.get("/api/web/notice_count").then(res => {
            this.noticeCount = res.data;
        });
    },

    methods: {
        setToken() {
            const token = localStorage.getItem("token");
            if (!token) {
                window.location.href = '/login';
                return
            }
            axios.get("/api/web/checkAuth", {
                headers: {"Authorization": "Bearer " + `${token}`},
            }).then(response => {
                console.log("User is authenticated.");
            }).catch(error => {
                console.log("User is not authenticated.");
                console.log(error);
                window.location.href = "/login";
            });
            axios.interceptors.request.use(config => {
                if (token) {
                    config.headers.Authorization = "Bearer " + `${token}`;
                }
                return config;
            }, err => {
                return Promise.reject(err);
            });
            return true;
        },
    },
    components: {
        'navi': {
            props: ['number'],
            template: `
<div id="nav">
    <div class="left">
        <a class="nav-item" href="/">
            <img class="icon" src="/static/img/home.png" alt="">
        </a>
        <a class="nav-item" href="/rule.html">
            <img class="icon" src="/static/img/rule.png" alt="">
        </a>
    </div>
    <div class="right">
        <a class="nav-item" href="/notice.html">
            <img class="icon" v-if="number==0" src="/static/img/notice_n.png" alt="">
            <img class="icon" v-else src="/static/img/notice.png" alt="">
        </a>
        <a class="nav-item" href="/setting.html">
            <img class="icon" src="/static/img/setting.png" alt="">
        </a>
    </div>
</div>
<div style="height: 4rem"></div>`
        }
    },
};
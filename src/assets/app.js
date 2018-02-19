/*const vm = new Vue({
    el: '#app',
    data: {
        results: [
            { title: "the very first post", abstract: "lorem ipsum some test dimpsum" },
            { title: "and then there was the second", abstract: "lorem ipsum some test dimsum" },
            { title: "third time's a charm", abstract: "lorem ipsum some test dimsum" },
            { title: "four the last time", abstract: "lorem ipsum some test dimsum" }
        ]
    }
});*/



const vm = new Vue({
    el: '#app',
    data: {
        results: [],
        total: 0,
        request: {
            /*
                       "requestor": {
                           "id": "ali",
                           "displayname": "ali",
                           "shortname": "hi",
                           "type": 0,
                           "domain": "mydom",
                           "publickeys": { "a": "b", "c": "D" }
                       },
                       "timestamp": 1518993431,*/
            "query": {
                "entry_type": 6,
                //"text": "الأندلس",
                //"tags": ["حِكمَة", "شِعر", "مَدح"],
                "limit": 20,
                "offset": 10
            }
        }
    },
    mounted() {
        axios.post("/api/entry/query", this.request)
            .then(response => {
                //console.log(response.data);
                this.total = response.data.total
                this.results = response.data.entries
            })
            .catch(function(error) {
                console.log(error);
            })
    },
    methods: {
        next: function() {
            if (this.request.query.offset + this.request.query.limit < this.total) {
                this.request.query.offset = this.request.query.offset + this.request.query.limit
                axios.post("/api/entry/query", this.request)
                    .then(response => {
                        //console.log(response.data);
                        this.total = response.data.total
                        this.results = response.data.entries
                    })
                    .catch(function(error) {
                        console.log(error);
                    })
            }
        },
        previous: function() {
            if (this.request.query.offset - this.request.query.limit >= 0) {
                this.request.query.offset = this.request.query.offset - this.request.query.limit
                axios.post("/api/entry/query", this.request)
                    .then(response => {
                        //console.log(response.data);
                        this.total = response.data.total
                        this.results = response.data.entries
                    })
                    .catch(function(error) {
                        console.log(error);
                    })
            }

        }
    }
});
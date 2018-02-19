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

query = {
    "requestor": {
        "id": "ali",
        "displayname": "ali",
        "shortname": "hi",
        "type": 0,
        "domain": "mydom",
        "publickeys": { "a": "b", "c": "D" }
    },
    "timestamp": 1518993431,
    "query": {
        "entry_type": 6,
        "limit": 5,
        "offset": 10
    }
}

const vm = new Vue({
    el: '#app',
    data: {
        results: [],
        total: 0
    },
    mounted() {
        axios.post("/api/entry/query", query)
            .then(response => {
                console.log(response.data);
                this.total = response.data.total
                this.results = response.data.entries
            })
            .catch(function(error) {
                console.log(error);
            })
    }
});
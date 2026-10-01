#include <climits>
#include <functional>
#include <queue>
#include <utility>
#include <vector>
using namespace std;

int delay(int n, vector<vector<int>>& edges, int source) {
    vector<vector<pair<int, int>>> adj(n);
    for (auto& e : edges) adj[e[0]].push_back({e[1], e[2]});
    vector<int> dist(n, INT_MAX);
    priority_queue<pair<int, int>, vector<pair<int, int>>, greater<>> pq;
    dist[source] = 0;
    pq.push({0, source});
    while (!pq.empty()) {
        auto [d, u] = pq.top();
        pq.pop();
        if (d > dist[u]) continue;
        for (auto [v, w] : adj[u])
            if (d + w < dist[v]) {
                dist[v] = d + w;
                pq.push({dist[v], v});
            }
    }
    int longest = 0;
    for (int d : dist) {
        if (d == INT_MAX) return -1;
        longest = max(longest, d);
    }
    return longest;
}

#include <climits>
#include <functional>
#include <queue>
#include <vector>
using namespace std;

int earliestArrival(int n, vector<int>& period, vector<vector<int>>& roads) {
    vector<vector<pair<int, int>>> out(n);
    for (auto& r : roads) out[r[0]].push_back({r[1], r[2]});
    vector<int> best(n, INT_MAX);
    best[0] = 0;
    priority_queue<pair<int, int>, vector<pair<int, int>>, greater<>> pq;
    pq.push({0, 0});
    while (!pq.empty()) {
        auto [d, u] = pq.top();
        pq.pop();
        if (d > best[u]) continue;
        int leave = (d + period[u] - 1) / period[u] * period[u];
        for (auto [v, t] : out[u])
            if (leave + t < best[v]) {
                best[v] = leave + t;
                pq.push({best[v], v});
            }
    }
    return best[n - 1] == INT_MAX ? -1 : best[n - 1];
}

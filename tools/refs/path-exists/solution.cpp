#include <vector>
using namespace std;

bool pathExists(int n, vector<vector<int>>& edges, int source, int target) {
    vector<vector<int>> adj(n);
    for (auto& e : edges) {
        adj[e[0]].push_back(e[1]);
        adj[e[1]].push_back(e[0]);
    }
    vector<bool> seen(n);
    vector<int> stack = {source};
    seen[source] = true;
    while (!stack.empty()) {
        int u = stack.back();
        stack.pop_back();
        if (u == target) return true;
        for (int v : adj[u])
            if (!seen[v]) {
                seen[v] = true;
                stack.push_back(v);
            }
    }
    return false;
}

#include <vector>
using namespace std;

bool twoColourable(int n, vector<vector<int>>& edges) {
    vector<vector<int>> adj(n);
    for (auto& e : edges) {
        adj[e[0]].push_back(e[1]);
        adj[e[1]].push_back(e[0]);
    }
    vector<int> colour(n, -1);
    for (int s = 0; s < n; s++) {
        if (colour[s] >= 0) continue;
        colour[s] = 0;
        vector<int> stack = {s};
        while (!stack.empty()) {
            int u = stack.back();
            stack.pop_back();
            for (int v : adj[u]) {
                if (colour[v] < 0) {
                    colour[v] = 1 - colour[u];
                    stack.push_back(v);
                } else if (colour[v] == colour[u]) {
                    return false;
                }
            }
        }
    }
    return true;
}

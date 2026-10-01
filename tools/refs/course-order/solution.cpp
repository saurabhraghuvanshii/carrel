#include <vector>
using namespace std;

// Depth-first: a course goes into the order after everything it leads to,
// then the order is reversed. Meeting a course still on the path means a loop.
static bool visit(int u, vector<vector<int>>& adj, vector<int>& state, vector<int>& out) {
    state[u] = 1;
    for (int v : adj[u]) {
        if (state[v] == 1) return false;
        if (state[v] == 0 && !visit(v, adj, state, out)) return false;
    }
    state[u] = 2;
    out.push_back(u);
    return true;
}

vector<int> order(int n, vector<vector<int>>& before) {
    vector<vector<int>> adj(n);
    for (auto& e : before) adj[e[0]].push_back(e[1]);
    vector<int> state(n), out;
    for (int u = 0; u < n; u++)
        if (state[u] == 0 && !visit(u, adj, state, out)) return {};
    return vector<int>(out.rbegin(), out.rend());
}

#include <algorithm>
#include <numeric>
#include <vector>
using namespace std;

// Kruskal: cheapest edges first, skipping any that would close a loop.
int cheapestConnection(int n, vector<vector<int>>& edges) {
    vector<vector<int>> e = edges;
    sort(e.begin(), e.end(), [](auto& a, auto& b) { return a[2] < b[2]; });
    vector<int> p(n);
    iota(p.begin(), p.end(), 0);
    auto root = [&](int x) {
        while (p[x] != x) x = p[x] = p[p[x]];
        return x;
    };
    int total = 0;
    for (auto& x : e) {
        int a = root(x[0]), b = root(x[1]);
        if (a != b) {
            p[a] = b;
            total += x[2];
        }
    }
    return total;
}

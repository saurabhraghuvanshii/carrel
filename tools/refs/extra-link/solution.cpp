#include <numeric>
#include <vector>
using namespace std;

vector<int> extraLink(int n, vector<vector<int>>& edges) {
    vector<int> p(n);
    iota(p.begin(), p.end(), 0);
    auto root = [&](int x) {
        while (p[x] != x) x = p[x] = p[p[x]];
        return x;
    };
    for (auto& e : edges) {
        int a = root(e[0]), b = root(e[1]);
        if (a == b) return e;
        p[a] = b;
    }
    return {};
}

#include <numeric>
#include <vector>
using namespace std;

static int findRoot(vector<int>& p, int x) {
    while (p[x] != x) x = p[x] = p[p[x]];
    return x;
}

int countGroups(int n, vector<vector<int>>& edges) {
    vector<int> p(n);
    iota(p.begin(), p.end(), 0);
    int groups = n;
    for (auto& e : edges) {
        int a = findRoot(p, e[0]), b = findRoot(p, e[1]);
        if (a != b) {
            p[a] = b;
            groups--;
        }
    }
    return groups;
}

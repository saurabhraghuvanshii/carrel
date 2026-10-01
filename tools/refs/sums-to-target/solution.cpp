#include <algorithm>
#include <vector>
using namespace std;

static void walk(const vector<int>& c, size_t from, int left, vector<int>& cur, vector<vector<int>>& out) {
    if (left == 0) {
        out.push_back(cur);
        return;
    }
    for (size_t i = from; i < c.size() && c[i] <= left; i++) {
        cur.push_back(c[i]);
        walk(c, i, left - c[i], cur, out);
        cur.pop_back();
    }
}

vector<vector<int>> combinations(vector<int>& candidates, int target) {
    vector<int> c = candidates;
    sort(c.begin(), c.end());
    vector<vector<int>> out;
    vector<int> cur;
    walk(c, 0, target, cur, out);
    return out;
}

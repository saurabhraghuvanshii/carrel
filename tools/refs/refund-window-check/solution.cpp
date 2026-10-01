#include <unordered_map>
#include <unordered_set>
#include <vector>
using namespace std;

vector<bool> approve(vector<vector<int>>& purchases, vector<vector<int>>& requests, int window) {
    unordered_map<int, int> bought;
    for (auto& p : purchases) bought[p[0]] = p[1];
    unordered_set<int> refunded;
    vector<bool> out;
    for (auto& r : requests) {
        auto it = bought.find(r[0]);
        bool ok = it != bought.end() && r[1] >= it->second && r[1] - it->second <= window && !refunded.count(r[0]);
        if (ok) refunded.insert(r[0]);
        out.push_back(ok);
    }
    return out;
}

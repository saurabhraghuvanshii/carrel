#include <map>
#include <utility>
#include <vector>
using namespace std;

vector<int> lookup(vector<vector<int>>& ranges, vector<int>& codes) {
    map<int, pair<int, int>> byStart;  // lo -> {hi, zone}
    for (auto& r : ranges) byStart[r[0]] = {r[1], r[2]};
    vector<int> out;
    for (int code : codes) {
        auto it = byStart.upper_bound(code);
        if (it == byStart.begin() || prev(it)->second.first < code) out.push_back(-1);
        else out.push_back(prev(it)->second.second);
    }
    return out;
}

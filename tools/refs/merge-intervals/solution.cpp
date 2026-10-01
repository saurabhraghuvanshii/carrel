#include <algorithm>
#include <vector>
using namespace std;

vector<vector<int>> merge(vector<vector<int>>& intervals) {
    vector<vector<int>> s = intervals, out;
    sort(s.begin(), s.end());
    for (auto& iv : s) {
        if (!out.empty() && iv[0] <= out.back()[1]) out.back()[1] = max(out.back()[1], iv[1]);
        else out.push_back(iv);
    }
    return out;
}

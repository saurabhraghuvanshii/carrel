#include <algorithm>
#include <vector>
using namespace std;

// Put every busy interval in one list, sort, and look for gaps.
vector<vector<int>> freeTime(vector<vector<vector<int>>>& busy) {
    vector<vector<int>> all;
    for (auto& person : busy)
        for (auto& iv : person) all.push_back(iv);
    sort(all.begin(), all.end());
    vector<vector<int>> free;
    int end = all[0][1];
    for (auto& iv : all) {
        if (iv[0] > end) free.push_back({end, iv[0]});
        end = max(end, iv[1]);
    }
    return free;
}

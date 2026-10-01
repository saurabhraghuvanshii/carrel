#include <algorithm>
#include <vector>
using namespace std;

vector<vector<int>> insert(vector<vector<int>>& intervals, vector<int>& add) {
    vector<vector<int>> out;
    size_t i = 0, n = intervals.size();
    while (i < n && intervals[i][1] < add[0]) out.push_back(intervals[i++]);
    int start = add[0], end = add[1];
    for (; i < n && intervals[i][0] <= end; i++) {
        start = min(start, intervals[i][0]);
        end = max(end, intervals[i][1]);
    }
    out.push_back({start, end});
    while (i < n) out.push_back(intervals[i++]);
    return out;
}

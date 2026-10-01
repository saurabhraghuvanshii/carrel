#include <algorithm>
#include <climits>
#include <vector>
using namespace std;

int fewestRemovals(vector<vector<int>>& intervals) {
    vector<vector<int>> s = intervals;
    sort(s.begin(), s.end(), [](auto& a, auto& b) { return a[1] < b[1]; });
    int kept = 0;
    long long last = LLONG_MIN;
    for (auto& iv : s) {
        if (iv[0] >= last) {
            kept++;
            last = iv[1];
        }
    }
    return (int)s.size() - kept;
}

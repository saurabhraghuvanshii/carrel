#include <algorithm>
#include <vector>
using namespace std;

// Sweep the sorted starts and ends: an end at the same time comes first.
int roomsNeeded(vector<vector<int>>& meetings) {
    vector<int> starts, ends;
    for (auto& m : meetings) {
        starts.push_back(m[0]);
        ends.push_back(m[1]);
    }
    sort(starts.begin(), starts.end());
    sort(ends.begin(), ends.end());
    int rooms = 0, best = 0;
    size_t e = 0;
    for (int s : starts) {
        while (e < ends.size() && ends[e] <= s) {
            e++;
            rooms--;
        }
        best = max(best, ++rooms);
    }
    return best;
}

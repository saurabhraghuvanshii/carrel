#include <algorithm>
#include <unordered_map>
#include <vector>
using namespace std;

vector<int> sessions(vector<vector<int>>& events, int timeout) {
    unordered_map<int, vector<int>> byUser;
    for (auto& e : events) byUser[e[0]].push_back(e[1]);
    int count = 0, longest = 0;
    for (auto& [user, times] : byUser) {
        sort(times.begin(), times.end());
        int start = times[0];
        count++;
        for (size_t i = 1; i < times.size(); i++) {
            if (times[i] - times[i - 1] > timeout) {
                count++;
                start = times[i];
            }
            longest = max(longest, times[i] - start);
        }
    }
    return {count, longest};
}

#include <unordered_map>
#include <vector>
using namespace std;

vector<int> funnel(vector<vector<int>>& events, int steps) {
    unordered_map<int, int> reached;
    vector<int> counts(steps, 0);
    for (auto& e : events) {
        int& r = reached[e[0]];
        if (e[1] == r + 1) {
            r++;
            counts[r - 1]++;
        }
    }
    return counts;
}

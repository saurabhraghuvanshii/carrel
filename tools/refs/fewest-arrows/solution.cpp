#include <algorithm>
#include <climits>
#include <vector>
using namespace std;

int fewestArrows(vector<vector<int>>& balloons) {
    vector<vector<int>> b = balloons;
    sort(b.begin(), b.end(), [](auto& x, auto& y) { return x[1] < y[1]; });
    int arrows = 0;
    long long at = LLONG_MIN;
    for (auto& balloon : b) {
        if (balloon[0] > at) {
            arrows++;
            at = balloon[1];
        }
    }
    return arrows;
}

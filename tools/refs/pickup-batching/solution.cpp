#include <map>
#include <utility>
#include <vector>
using namespace std;

static int floorDiv(int a, int b) {
    return a / b - (a % b != 0 && a < 0);
}

int countBatches(vector<vector<int>>& requests, int limit, int window) {
    map<pair<int, int>, pair<int, int>> open;
    int count = 0;
    for (auto& r : requests) {
        auto zone = make_pair(floorDiv(r[1], 100), floorDiv(r[2], 100));
        auto it = open.find(zone);
        if (it != open.end() && r[0] - it->second.first <= window && it->second.second < limit) {
            it->second.second++;
        } else {
            open[zone] = {r[0], 1};
            count++;
        }
    }
    return count;
}

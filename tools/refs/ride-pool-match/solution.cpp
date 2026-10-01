#include <algorithm>
#include <vector>
using namespace std;

int maxRiders(int stops, int seats, vector<vector<int>>& rides) {
    sort(rides.begin(), rides.end(), [](auto& a, auto& b) { return a[1] < b[1]; });
    vector<int> load(stops, 0);
    int taken = 0;
    for (auto& r : rides) {
        if (*max_element(load.begin() + r[0], load.begin() + r[1]) >= seats) continue;
        for (int x = r[0]; x < r[1]; x++) load[x]++;
        taken++;
    }
    return taken;
}

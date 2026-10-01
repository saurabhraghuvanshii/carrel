#include <climits>
#include <cstdlib>
#include <vector>
using namespace std;

int nearestDriver(vector<vector<int>>& drivers, int rx, int ry) {
    int best = -1;
    long long bestDist = LLONG_MAX;
    for (auto& d : drivers) {
        if (d[3] == 0) continue;
        long long dist = llabs((long long)d[1] - rx) + llabs((long long)d[2] - ry);
        if (dist < bestDist || (dist == bestDist && d[0] < best)) {
            best = d[0];
            bestDist = dist;
        }
    }
    return best;
}

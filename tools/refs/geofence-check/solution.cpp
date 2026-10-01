#include <vector>
using namespace std;

int entries(vector<int>& zone, vector<vector<int>>& path) {
    int count = 0;
    bool wasInside = false;
    for (auto& p : path) {
        bool inside = zone[0] <= p[0] && p[0] <= zone[2] && zone[1] <= p[1] && p[1] <= zone[3];
        if (inside && !wasInside) count++;
        wasInside = inside;
    }
    return count;
}

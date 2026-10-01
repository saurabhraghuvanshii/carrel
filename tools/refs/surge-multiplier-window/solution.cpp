#include <algorithm>
#include <vector>
using namespace std;

vector<int> surgeMultipliers(vector<int>& requests, vector<int>& drivers, int k) {
    vector<int> out;
    int r = 0, d = 0;
    for (int i = 0; i < (int)requests.size(); i++) {
        r += requests[i];
        d += drivers[i];
        if (i >= k) r -= requests[i - k], d -= drivers[i - k];
        if (i >= k - 1) out.push_back(r <= d ? 1 : d == 0 ? 5 : min(5, (r + d - 1) / d));
    }
    return out;
}

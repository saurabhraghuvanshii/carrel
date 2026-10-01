#include <algorithm>
#include <climits>
#include <vector>
using namespace std;

vector<int> nearest(vector<int>& stores, vector<int>& customers) {
    sort(stores.begin(), stores.end());
    vector<int> out;
    for (int c : customers) {
        auto it = lower_bound(stores.begin(), stores.end(), c);
        int best = INT_MAX;
        if (it != stores.end()) best = *it - c;
        if (it != stores.begin()) best = min(best, c - *prev(it));
        out.push_back(best);
    }
    return out;
}

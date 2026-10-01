#include <algorithm>
#include <numeric>
#include <vector>
using namespace std;

int maxProfit(vector<vector<int>>& jobs) {
    sort(jobs.begin(), jobs.end(), [](auto& x, auto& y) { return x[1] > y[1]; });
    vector<int> latest(jobs.size() + 1);
    iota(latest.begin(), latest.end(), 0);
    auto find = [&](int t) {
        while (latest[t] != t) t = latest[t] = latest[latest[t]];
        return t;
    };
    int total = 0;
    for (auto& j : jobs) {
        int t = find(j[0]);
        if (t > 0) {
            total += j[1];
            latest[t] = t - 1;
        }
    }
    return total;
}

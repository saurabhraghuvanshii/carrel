#include <algorithm>
#include <vector>
using namespace std;

long long finalTotal(vector<int>& prices, vector<vector<int>>& coupons) {
    vector<vector<int>> active;
    for (auto& c : coupons) {
        if (c[0] != 0) active.push_back(c);
        else if (!active.empty()) active.pop_back();
    }
    long long total = 0;
    for (int p : prices) total += p;
    for (auto& c : active) {
        if (c[0] == 1) total -= total * c[1] / 100;
        else total = max(0LL, total - c[1]);
    }
    return total;
}

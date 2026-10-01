#include <algorithm>
#include <climits>
#include <vector>
using namespace std;

int longestChain(vector<vector<int>>& pairs) {
    sort(pairs.begin(), pairs.end(), [](auto& x, auto& y) { return x[1] < y[1]; });
    int chain = 0;
    long long end = LLONG_MIN;
    for (auto& p : pairs)
        if (p[0] > end) {
            chain++;
            end = p[1];
        }
    return chain;
}

#include <climits>
#include <vector>
using namespace std;

int cheapestPrice(int n, vector<vector<int>>& flights, int from, int to, int k) {
    vector<long long> cost(n, LLONG_MAX);
    cost[from] = 0;
    for (int round = 0; round <= k; round++) {
        vector<long long> next = cost;
        for (auto& f : flights)
            if (cost[f[0]] != LLONG_MAX && cost[f[0]] + f[2] < next[f[1]]) next[f[1]] = cost[f[0]] + f[2];
        cost = next;
    }
    return cost[to] == LLONG_MAX ? -1 : (int)cost[to];
}

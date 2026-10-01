#include <algorithm>
#include <cstdlib>
#include <vector>
using namespace std;

long long minTotalDistance(vector<int>& orders, vector<int>& couriers) {
    sort(orders.begin(), orders.end());
    sort(couriers.begin(), couriers.end());
    int n = orders.size(), m = couriers.size();
    const long long inf = 1LL << 60;
    vector<vector<long long>> best(n + 1, vector<long long>(m + 1, inf));
    for (int j = 0; j <= m; j++) best[0][j] = 0;
    for (int i = 1; i <= n; i++)
        for (int j = i; j <= m; j++)
            best[i][j] = min(best[i][j - 1], best[i - 1][j - 1] + abs(orders[i - 1] - couriers[j - 1]));
    return best[n][m];
}

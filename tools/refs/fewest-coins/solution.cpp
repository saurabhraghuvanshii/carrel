#include <algorithm>
#include <vector>
using namespace std;

int fewestCoins(vector<int>& coins, int amount) {
    vector<int> best(amount + 1, -1);
    best[0] = 0;
    for (int v = 1; v <= amount; v++)
        for (int c : coins)
            if (c <= v && best[v - c] >= 0 && (best[v] < 0 || best[v - c] + 1 < best[v])) best[v] = best[v - c] + 1;
    return best[amount];
}

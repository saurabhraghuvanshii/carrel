#include <vector>
using namespace std;

int countWays(vector<int>& coins, int amount) {
    vector<long long> ways(amount + 1, 0);
    ways[0] = 1;
    for (int c : coins)
        for (int v = c; v <= amount; v++) ways[v] += ways[v - c];
    return (int)ways[amount];
}

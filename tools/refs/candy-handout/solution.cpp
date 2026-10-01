#include <algorithm>
#include <vector>
using namespace std;

int fewestSweets(vector<int>& ratings) {
    int n = ratings.size();
    vector<int> sweets(n, 1);
    for (int i = 1; i < n; i++)
        if (ratings[i] > ratings[i - 1]) sweets[i] = sweets[i - 1] + 1;
    int total = sweets[n - 1];
    for (int i = n - 2; i >= 0; i--) {
        if (ratings[i] > ratings[i + 1]) sweets[i] = max(sweets[i], sweets[i + 1] + 1);
        total += sweets[i];
    }
    return total;
}

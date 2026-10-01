#include <algorithm>
#include <vector>
using namespace std;

int bestValue(vector<int>& weights, vector<int>& values, int capacity) {
    vector<int> best(capacity + 1, 0);
    for (size_t i = 0; i < weights.size(); i++)
        for (int c = capacity; c >= weights[i]; c--) best[c] = max(best[c], best[c - weights[i]] + values[i]);
    return best[capacity];
}

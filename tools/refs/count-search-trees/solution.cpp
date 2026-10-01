#include <vector>
using namespace std;

int countTrees(int n) {
    vector<long long> trees(n + 1, 0);
    trees[0] = 1;
    for (int k = 1; k <= n; k++)
        for (int root = 1; root <= k; root++) trees[k] += trees[root - 1] * trees[k - root];
    return (int)trees[n];
}

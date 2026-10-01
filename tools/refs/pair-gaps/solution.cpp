#include <algorithm>
#include <cstdlib>
#include <vector>
using namespace std;

long long smallestGapSum(vector<int>& a, vector<int>& b) {
    sort(a.begin(), a.end());
    sort(b.begin(), b.end());
    long long total = 0;
    for (size_t i = 0; i < a.size(); i++) total += abs(a[i] - b[i]);
    return total;
}

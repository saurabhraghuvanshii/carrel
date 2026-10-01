#include <algorithm>
#include <climits>
#include <vector>
using namespace std;

double median(vector<int>& first, vector<int>& second) {
    if (first.size() > second.size()) return median(second, first);
    int m = first.size(), n = second.size(), half = (m + n + 1) / 2;
    int lo = 0, hi = m;
    while (true) {
        int i = (lo + hi) / 2, j = half - i;
        int leftA = i > 0 ? first[i - 1] : INT_MIN, rightA = i < m ? first[i] : INT_MAX;
        int leftB = j > 0 ? second[j - 1] : INT_MIN, rightB = j < n ? second[j] : INT_MAX;
        if (leftA <= rightB && leftB <= rightA) {
            int leftMax = max(leftA, leftB);
            if ((m + n) % 2 == 1) return leftMax;
            return (leftMax + (long long)min(rightA, rightB)) / 2.0;
        }
        if (leftA > rightB) hi = i - 1;
        else lo = i + 1;
    }
}

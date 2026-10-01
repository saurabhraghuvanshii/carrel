#include <algorithm>
#include <vector>
using namespace std;

int largestRectangle(vector<int>& heights) {
    int n = heights.size(), best = 0;
    vector<int> rising;
    for (int i = 0; i <= n; i++) {
        int cur = i < n ? heights[i] : 0;
        while (!rising.empty() && heights[rising.back()] >= cur) {
            int h = heights[rising.back()];
            rising.pop_back();
            int left = rising.empty() ? -1 : rising.back();
            best = max(best, h * (i - left - 1));
        }
        rising.push_back(i);
    }
    return best;
}

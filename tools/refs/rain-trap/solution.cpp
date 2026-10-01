#include <algorithm>
#include <vector>
using namespace std;

int trappedWater(vector<int>& heights) {
    int i = 0, j = (int)heights.size() - 1, leftMax = 0, rightMax = 0, total = 0;
    while (i <= j) {
        if (leftMax <= rightMax) {
            leftMax = max(leftMax, heights[i]);
            total += leftMax - heights[i++];
        } else {
            rightMax = max(rightMax, heights[j]);
            total += rightMax - heights[j--];
        }
    }
    return total;
}

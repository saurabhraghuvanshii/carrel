#include <vector>
using namespace std;

int missingNumber(vector<int>& nums) {
    long long n = nums.size(), sum = 0;
    for (int v : nums) sum += v;
    return n * (n + 1) / 2 - sum;
}

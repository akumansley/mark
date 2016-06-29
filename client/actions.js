const uid = () => Math.random().toString(34).slice(2);

export function addMark(url) {
  return {
    type: 'ADD_MARK',
    payload: {
      id: uid(),
      url: url
    }
  };
}

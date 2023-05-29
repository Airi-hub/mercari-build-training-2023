import React, { useEffect, useState } from 'react';

interface Item {
  id: number;
  name: string;
  category: string;
  image_filename: string;
}

const server = process.env.REACT_APP_API_URL || 'http://127.0.0.1:9000';
const placeholderImage = process.env.PUBLIC_URL + '/default.jpg';

interface Props {
  reload?: boolean;
  onLoadCompleted?: () => void;
}
// Propsの名前を修正
export const ItemList: React.FC<Props> = (props) => {
  const { reload = true, onLoadCompleted } = props;
  const [items, setItems] = useState<Item[]>([]);

  const fetchItems = () => {
    fetch(`${server}/items`, {
      method: 'GET',
      mode: 'cors',
      headers: {
        'Content-Type': 'application/json',
        'Accept': 'application/json'
      }
    })
    //レスポンスデータの取り扱いを修正
      .then(response => {
        if (response.ok) {
          return response.json();
        } else {
          throw new Error('GET error: ' + response.status);
        }
      })
      .then(data => {
        console.log('GET success:', data);
        setItems(data.items);
        onLoadCompleted && onLoadCompleted();
      })
      .catch(error => {
        console.error('GET error:', error);
      });
  };

  useEffect(() => {
    if (reload) {
      fetchItems();
    }
  }, [reload, onLoadCompleted]);  //依存配列にonLoadCompletedを追加

  return (
    <div className='ItemListing'>
      {items.map((item) => (
        <div key={item.id} className='ItemList'>
          <img
            src={`${server}/image/${item.image_filename}`}
            alt={item.name}
            onError={(e) => {
              e.currentTarget.src = placeholderImage;
            }}
          />
    
          <p>
            <span>Name: {item.name}</span>
            <br />
            <span>Category: {item.category}</span>
          </p>
        </div>
      ))}
    </div>
  );  
};
